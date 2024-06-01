resource "aws_security_group" "endpoint_access" {
  description = "Access to endpoints"
  name        = "${local.prefix}-endpoint-access"
  vpc_id      = aws_vpc.main.id

  ingress {
    cidr_blocks = [aws_vpc.main.cidr_block]
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
  }
}


resource "aws_security_group" "rds" {
  description = "Allow access to RDS database instance"
  name        = "${local.prefix}-rds-inbound-access"
  vpc_id      = aws_vpc.main.id

  ingress {
    protocol  = "tcp"
    from_port = 5432
    to_port   = 5432
  }

  tags = {
    Name = "${local.prefix}-db-security-group"
  }
}
