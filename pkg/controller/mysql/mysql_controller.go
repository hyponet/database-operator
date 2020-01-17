package mysql

import (
	"context"
	databasev1 "database-operator/pkg/apis/database/v1"
	"database-operator/pkg/mysql"
	"database-operator/pkg/utils/stsutil"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_mysql")

const controllerName = "controller-mysql"

// Add creates a new MySQL Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMySQL{
		client:      mgr.GetClient(),
		scheme:      mgr.GetScheme(),
		eventRecord: mgr.GetEventRecorderFor(controllerName),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MySQL
	err = c.Watch(&source.Kind{Type: &databasev1.MySQL{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource StatefulSets and requeue the owner MySQL
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &databasev1.MySQL{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileMySQL implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMySQL{}

// ReconcileMySQL reconciles a MySQL object
type ReconcileMySQL struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client      client.Client
	eventRecord record.EventRecorder
	scheme      *runtime.Scheme
}

// Reconcile reads that state of the cluster for a MySQL object and makes changes based on the state read
// and what is in the MySQL.Spec
func (r *ReconcileMySQL) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.TODO()
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MySQL")

	// Fetch the MySQL instance
	instance := &databasev1.MySQL{}
	err := r.client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new mysql instance
	mysqlInstance := mysql.NewMySqlInstance(instance.DeepCopy())
	sts := mysqlInstance.NewStsForCR()

	// Set MySQL instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, sts, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this StatefulSet already exists
	found := &appsv1.StatefulSet{}
	err = r.client.Get(ctx, types.NamespacedName{Name: sts.Name, Namespace: sts.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new StatefulSet", "sts.Namespace", sts.Namespace, "sts.Name", sts.Name)
		err = r.client.Create(ctx, sts)
		if err != nil {
			return reconcile.Result{}, err
		}

		// StatefulSet created successfully - don't requeue
		r.recordMySqlInstanceEvent(instance, stsCreated, fmt.Sprintf("Create new StatefulSet %s", sts.Name))
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	status, err := mysqlInstance.UpdateStatus(found.DeepCopy())
	if err != nil {
		reqLogger.Error(err, "Get MysqlInstance status error")
	}
	if status == nil {
		// nothing changed - don't requeue
		return reconcile.Result{}, nil
	}

	if err := r.updateMySqlInstanceStatusWithEventRecord(instance, status); err != nil {
		reqLogger.Error(err, "Update MysqlInstance status error")
	}

	if !stsutil.DeepEqual(sts, found) {
		reqLogger.Info("Update Old StatefulSet", "sts.Namespace", found.Namespace, "sts.Name", found.Name)
		err = r.client.Update(ctx, sts)
		if err != nil {
			reqLogger.Error(err, "Update StatefulSet error", "sts.Namespace", found.Namespace, "sts.Name", found.Name)
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileMySQL) updateMySqlInstanceStatusWithEventRecord(instance *databasev1.MySQL, status *databasev1.MySQLStatus) error {
	return nil
}
