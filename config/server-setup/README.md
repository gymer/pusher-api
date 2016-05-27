# rails-backend-playbook
Ansible playbook for setup Rails 4 application environment: ruby-2.1.5, deploy user, nginx, postgresql-9.3, nodejs, unicorn.

## Config

```
cp vars/defaults.example.yml vars/defaults.yml
```
then in `vars/defaults.yml` set `app_name`, `env` and `webserver_name` params.
Your dont need keep `config/database.yml` and `config/secrets.yml` in VCS, it will be generated automatically for `env` that you have defined.

## Setup locally with Vagrant

```
vagrant up
```

then you can check out that all dependencies installed:

```
ssh deploy@127.0.0.1 -p 2200
ruby -v
node -v
psql --version
psql app_{{app_name}}_{{env}}
\q
exit
```

## Setup remote server

Create inventory file `hosts` and run:

```
ansible-playbook -i hosts playbook.yml -e "secret=SECRET_KEY_BASE"
```

## Deploy with capistrano 3

In your rails app directory add to Gemfile:

```
gem 'capistrano', '~> 3.4.0'
gem 'capistrano-rails',   require: false
gem 'capistrano-bundler', require: false
gem 'capistrano-unicorn-nginx', '~> 3.4.0'
```

run

```
bundle install
cap install
```

add to  `Capfile`

```
require 'capistrano/rails'
require 'capistrano/bundler'
require 'capistrano/unicorn_nginx'
```

open `config/deploy/production.rb` and define in  params (for Vagrant set server `127.0.0.1` and port `2222`):

```
server 'your-server-address', user: 'deploy', port: 22, roles: [:web, :app, :db], primary: true
set :branch, ENV['BRANCH'] || 'your-production-branch-name'
```

copy deploy config from `rails-ansible-playbook`
```
cp config/deploy.rb /your/rails-app/config/deploy.rb
```

then run

```
cap setup
cap production deploy:check
cap production deploy
```
