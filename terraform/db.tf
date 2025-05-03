resource "aws_rds_cluster" "mst_db" {
  cluster_identifier     = "mst-aurora-db"
  engine                 = "aurora-postgresql"
  engine_mode            = "provisioned"
  engine_version         = "14.12"
  database_name          = "mst_db"
  master_username        = "mst_admin"
  master_password        = var.db_password
  storage_encrypted      = true
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.mst_db_sg.id]

  skip_final_snapshot       = true
  final_snapshot_identifier = "mst-final-${random_id.snapshot_suffix.hex}"
  apply_immediately         = true

  serverlessv2_scaling_configuration {
    max_capacity             = 1.0
    min_capacity             = 0.0
    seconds_until_auto_pause = 3600
  }
}
resource "random_id" "snapshot_suffix" {
  byte_length = 4
}

resource "aws_rds_cluster_instance" "cluster" {
  cluster_identifier  = aws_rds_cluster.mst_db.id
  publicly_accessible = true
  instance_class      = "db.serverless"
  engine              = aws_rds_cluster.mst_db.engine
  engine_version      = aws_rds_cluster.mst_db.engine_version
}
resource "aws_db_subnet_group" "main" {
  name       = "mst-db-subnet-group"
  subnet_ids = [aws_subnet.public_c.id, aws_subnet.public_b.id]

  tags = {
    Name = "MST DB Subnet Group"
  }
}


resource "aws_security_group" "mst_db_sg" {
  name        = "mst_db_sg"
  description = "Security group for MST Aurora DB"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port = 5432
    to_port   = 5432
    protocol  = "tcp"
    //TODO: Change this in production to just my IP's
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
