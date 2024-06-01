######################################################################
# Endpoint to allow ECS to access ECR, CloudWatch and System Manager #
######################################################################

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

resource "aws_vpc_endpoint" "ecr" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.ecr.api"
  vpc_endpoint_type = "Interface"

  subnet_ids = [aws_subnet.private_a.id, aws_subnet.private_b.id]

  security_group_ids = [
    aws_security_group.endpoint_access.id
  ]

  tags = {
    Name = "${local.prefix}-ecr-endpoint"
  }
}

resource "aws_vpc_endpoint" "dkr" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.ecr.dkr"
  vpc_endpoint_type = "Interface"

  subnet_ids = [aws_subnet.private_a.id, aws_subnet.private_b.id]

  security_group_ids = [
    aws_security_group.endpoint_access.id
  ]

  tags = {
    Name = "${local.prefix}-dkr-endpoint"
  }
}

resource "aws_vpc_endpoint" "cloudwatch_logs" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.logs"
  vpc_endpoint_type = "Interface"

  subnet_ids = [aws_subnet.private_a.id, aws_subnet.private_b.id]

  security_group_ids = [
    aws_security_group.endpoint_access.id
  ]

  tags = {
    Name = "${local.prefix}-cloudwatch-endpoint"
  }
}

resource "aws_vpc_endpoint" "ssm" {
  vpc_id              = aws_vpc.main.id
  service_name        = "com.amazonaws.${data.aws_region.current.name}.ssmmessages"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true

  subnet_ids = [aws_subnet.private_a.id, aws_subnet.private_b.id]

  security_group_ids = [
    aws_security_group.endpoint_access.id
  ]

  tags = {
    Name = "${local.prefix}-ssmmessages-endpoint"
  }
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.s3"
  vpc_endpoint_type = "Gateway"

  route_table_ids = [
    aws_vpc.main.default_route_table_id
  ]

  tags = {
    Name = "${local.prefix}-s3-endpoint"
  }
}
