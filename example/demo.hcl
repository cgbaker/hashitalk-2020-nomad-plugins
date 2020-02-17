job "demo" {
  datacenters = ["dc1"]
  type = "service"
  group "demo" {
    task "python-server" {

      driver = "python"
      config {
        script = "server.py"
      }

      template {
        destination = "server.py"
        data = <<EOF
import SimpleHTTPServer
import SocketServer
import os

PORT = int(os.getenv("NOMAD_PORT_http", "8080"))
ADDR = os.getenv("NOMAD_IP_http", "localhost")

web_dir = os.path.join(os.path.dirname(__file__), 'www')
os.chdir(web_dir)

Handler = SimpleHTTPServer.SimpleHTTPRequestHandler

httpd = SocketServer.TCPServer((ADDR, PORT), Handler)

print "serving at port", PORT
httpd.serve_forever()
EOF
      }

      template {
        destination = "www/index.html"
        data = <<EOF
          <html>
          <h1>HashiTalks 2020: Extending Nomad with Plugins</h1>
          <p> Source code
          <ul>
          <li><a href="https://github.com/cgbaker/hashitalk-2020-nomad-plugins">https://github.com/cgbaker/hashitalk-2020-nomad-plugins</a>
          </ul>
          </html>
EOF
      }

      resources {
        cpu = 512
        memory = 128
        network {
          mbits = 100
          port "http" {}
        }
      }
     
    }
  }
}
