job "device-demo" {
  datacenters = ["dc1"]
  type = "service"
  group "demo" {
    task "demo" {

      driver = "raw_exec"
      config {
        command = "sleep"
        args = ["120"]
      }

      resources {
        cpu = 512
        memory = 128
        device "usb" {}
      }
     
    }
  }
}
