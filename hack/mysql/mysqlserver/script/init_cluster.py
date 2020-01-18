import os
import socket

ENV = os.environ
HOSTNAME = socket.gethostname()
POD_IP = ENV.get("PODIP")
DB_USER = "root"
DB_PASSWD = os.getenv("MYSQL_ROOT_PASSWORD")

splits = HOSTNAME.split("-")
try:
    index = int(splits[-1])
except TypeError:
    raise TypeError("Get index error")

service_name = "-".join(splits[:-1])
cluster_name = service_name.replace("-", "")

if index == 0:
    shell.connect('{db_user}@localhost:3306'.format(db_user=DB_USER), DB_PASSWD)
    try:
        cluster = dba.get_cluster(cluster_name)
    except Exception:
        cluster = None
        dba.configure_local_instance(
            '{db_user}@localhost:3306'.format(db_user=DB_USER),
            {'password': DB_PASSWD, 'interactive': False},
        )
        dba.create_cluster(cluster_name, {'localAddress': POD_IP})

    if cluster is None:
        cluster = dba.get_cluster(cluster_name)
    cluster.rescan()

else:

    shell.connect('{db_user}@localhost:3306'.format(db_user=DB_USER), DB_PASSWD)
    try:
        dba.configure_local_instance(
            '{db_user}@localhost:3306'.format(db_user=DB_USER),
            {'password': DB_PASSWD, 'interactive': False},
        )
    except Exception:
        pass

    session.close()
    shell.connect('{db_user}@{service_name}-0.{service_name}:3306'.format(
        db_user=DB_USER,
        service_name=service_name,
    ), DB_PASSWD)

    cluster = dba.get_cluster(cluster_name)
    cluster.rescan()

    dbKey = '{service_name}-0.{service_name}:3306'.format(service_name=service_name)

    status = cluster.status()
    defaultReplicaSet = status.get('defaultReplicaSet', {})
    topology = defaultReplicaSet.get('defaultReplicaSet', {})

    if dbKey not in topology:
        cluster.add_instance('{db_user}@{pod_id}:3306'.format(db_user=DB_USER, pod_id=POD_IP),
                             {'localAddress': POD_IP, 'password': DB_PASSWD, 'recoveryMethod': 'auto'})
