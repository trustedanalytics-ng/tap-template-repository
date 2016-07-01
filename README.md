### REQUIREMENTS

GO 1.6

### LOCAL FROM SOURCES INSTALLATION

```
  git clone https://github.com/intel-data/tap-template-repository    
  ./pack.sh    
  TEMPLATE_REPOSITORY_USER=admin TEMPLATE_REPOSITORY_PASS=password TEMPLATE_REPOSITORY_PORT=8082 ./temp/bin/tap-template-repository
```

validate endpont


```
  curl -v  admin:password@localhost:8082/api/v1/templates
```