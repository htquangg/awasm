resource "aws_db_subnet_group" "main" {
  name = "${local.prefix}-main"
  subnet_ids = [
    aws_subnet.private_a.id,
    aws_subnet.private_b.id,
  ]

  tags = {
    Name = "${local.prefix}-db-subnet-group"
  }
}

resource "aws_db_instance" "main" {
  identifier                 = "${local.prefix}-db"
  db_name                    = "awasm"
  allocated_storage          = 20
  storage_type               = "gp2"
  engine                     = "postgres"
  engine_version             = "14.12"
  auto_minor_version_upgrade = true
  instance_class             = "db.t4g.micro"
  username                   = var.db_username
  password                   = var.db_password
  skip_final_snapshot        = true
  db_subnet_group_name       = aws_db_subnet_group.main.name
  multi_az                   = false
  backup_retention_period    = 0
  vpc_security_group_ids     = [aws_security_group.rds.id]

  tags = {
    Name = "${local.prefix}-main"
  }
}
