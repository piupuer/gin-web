# nginx配置反向代理去除端口

## 安装nginx(这里以Ubuntu系统为例)

### 1. 安装

```bash
apt-get install nginx

nginx -v
# 我目前安装的是
# nginx version: nginx/1.14.0 (Ubuntu)

```

### 2. 启动

```bash
systemctl start nginx
```

## 配置nginx.conf

### 1. 找到nginx.conf位置

```bash
nginx -t
# 输出如下
# nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
# nginx: configuration file /etc/nginx/nginx.conf test is successful
```

### 2. 替换nginx.conf

```bash
worker_processes auto;
pid /run/nginx.pid;

events {
  worker_connections 768;
  # multi_accept on;
}

http {

  ##
  # Basic Settings
  ##
  autoindex on;
  autoindex_exact_size off;
  client_max_body_size 5m;
  
  sendfile on;
  tcp_nopush on;
  tcp_nodelay on;
  keepalive_timeout 65;
  types_hash_max_size 2048;
  # server_tokens off;

  # server_names_hash_bucket_size 64;
  # server_name_in_redirect off;

  include /etc/nginx/mime.types;
  default_type application/octet-stream;

  ##
  # SSL Settings
  ##

  ssl_protocols TLSv1 TLSv1.1 TLSv1.2; # Dropping SSLv3, ref: POODLE
  ssl_prefer_server_ciphers on;

  ##
  # Logging Settings
  ##

  access_log /var/log/nginx/access.log;
  error_log /var/log/nginx/error.log;

  ##
  # Gzip Settings
  ##

  gzip on;
  gzip_disable "msie6";

  # gzip_vary on;
  # gzip_proxied any;
  # gzip_comp_level 6;
  # gzip_buffers 16 8k;
  # gzip_http_version 1.1;
  # gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

  ##
  # Virtual Host Configs
  ##

  # 引入虚拟主机配置
  include /etc/nginx/conf.d/*.conf;
  include /etc/nginx/sites-enabled/*;
}


#mail {
#  # See sample authentication script at:
#  # http://wiki.nginx.org/ImapAuthenticateWithApachePhpScript
# 
#  # auth_http localhost/auth.php;
#  # pop3_capabilities "TOP" "USER";
#  # imap_capabilities "IMAP4rev1" "UIDPLUS";
# 
#  server {
#  listen    localhost:110;
#  protocol    pop3;
#  proxy    on;
#  }
# 
#  server {
#  listen    localhost:143;
#  protocol    imap;
#  proxy    on;
#  }
#}
```

> 特别注意这两行不能少  
include /etc/nginx/conf.d/*.conf;  
include /etc/nginx/sites-enabled/*;

### 3. 配置前后端虚拟主机映射

```bash
vim /etc/nginx/conf.d/gin-web.conf 
```

写入如下内容

```bash
# 后端应用映射
upstream gin-web {
  server 127.0.0.1:10000;
  keepalive 64;
}
# 前端应用映射
upstream gin-web-vue {
  server 127.0.0.1:10001;
  keepalive 64;
}
# pprof应用映射
upstream gin-pprof {
  server 127.0.0.1:10005;
  keepalive 64;
}
server {
  listen 80;
  # 开启https
  #listen 443 ssl;
  # 证书所在目录
  #ssl_certificate cert/domain.com.pem;
  #ssl_certificate_key cert/domain.com.key;
  # http自动重定向到https
  #if ( $ssl_protocol = "") {
  #  rewrite ^ https://$host$request_uri? permanent;
  #}        

  # 绑定域名
  # server_name domain.com;
  server_name 127.0.0.1;

  location / {
    proxy_redirect     off;
    proxy_set_header    X-Real-IP  $remote_addr;
    proxy_set_header    X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header    X-Forwarded-Proto $scheme;
    proxy_set_header    Host $http_host;
    proxy_set_header    X-NginX-Proxy true;
    proxy_set_header    Connection "";
    proxy_http_version 1.1;
    # 末尾加斜杠将不会转发location path(有二级目录时有用处)
    proxy_pass         http://gin-web-vue/;
  }
  
  location ^~ /api {
    proxy_redirect     off;
    proxy_set_header    X-Real-IP $remote_addr;
    proxy_set_header    X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header    X-Forwarded-Proto $scheme;
    proxy_set_header    Host $http_host;
    proxy_set_header    X-NginX-Proxy true;
    proxy_set_header    Upgrade $http_upgrade;
    proxy_set_header    Connection 'upgrade';
    proxy_http_version 1.1;
    proxy_pass         http://gin-web/api;
  }

  location ^~ /debug/pprof/ {
    proxy_redirect     off;
    proxy_set_header    X-Real-IP $remote_addr;
    proxy_set_header    X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header    X-Forwarded-Proto $scheme;
    proxy_set_header    Host $http_host;
    proxy_set_header    X-NginX-Proxy true;
    proxy_set_header    Connection "";
    proxy_http_version 1.1;
    proxy_pass         http://gin-pprof/debug/pprof/;
  }
}
```

### 4. 重启

```bash
systemctl restart nginx
```

### 5. 访问

现已成功去除端口, 浏览器输入127.0.0.1即可访问
