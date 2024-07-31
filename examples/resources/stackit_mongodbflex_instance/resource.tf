resource "stackit_mongodbflex_instance" "example" {
  project_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "example-instance"
  acl        = ["XXX.XXX.XXX.X/XX", "XX.XXX.XX.X/XX"]
  flavor = {
    cpu = 1
    ram = 8
  }
  replicas = 1
  storage = {
    class = "class"
    size  = 10
  }
  version = "7.0"
  options = {
    type                    = "Single"
    snapshot_retention_days = 3
  }
  backup_schedule = "0 0 * * *"
}
