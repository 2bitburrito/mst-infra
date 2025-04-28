resource "aws_rds_cluster" "mst_db" {
  cluster_identifier = "mst-aurora-db"
  engine             = "aurora-postgresql"
  engine_mode        = "provisioned"
  engine_version     = "14.12"
  database_name      = "mst_db"
  master_username    = "mst_admin"
  master_password    = var.db_password
  storage_encrypted  = true

  serverlessv2_scaling_configuration {
    max_capacity             = 1.0
    min_capacity             = 0.0
    seconds_until_auto_pause = 3600
  }
}

resource "aws_rds_cluster_instance" "cluster" {
  cluster_identifier  = aws_rds_cluster.mst_db.id
  publicly_accessible = true
  instance_class      = "db.serverless"
  engine              = aws_rds_cluster.mst_db.engine
  engine_version      = aws_rds_cluster.mst_db.engine_version
}
