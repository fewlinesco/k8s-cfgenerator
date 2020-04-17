# Usage

```
$> cfgenerator volumes/secrets volumes/config < config.jsonnet
{
   "api": {
      "address": "0.0.0.0:1337"
   },
   "database": {
      "password": "sssh! it's secret",
      "username": "myapp"
   }
}
```
