# Local deployment

```
git clone https://github.com/vz10/geogit.git
```

If you use Kitematic then run
```
# cd geogit
# eval "$(docker-machine env default)"
# docker build -t geogit .
```

If everything went ok, you can run `docker-compose up`
