# ngxcpd

nginx cache purge daemon (and inspector)

(Work in progress)

## roadmap

- [x] nginx cache file parser
- [x] scan cache directory to in-memory index
- [x] watch cache directory for changes
- [x] manually clear cache entries in a (local) directory using a pattern on cache key
- [ ] cache invalidation protocol (merge with emgag/varnish-towncrier?)
- [ ] pub/sub channel listen (different providers? redis, sns, azure service bus)
- [ ] config file (viper?)
- [ ] manually clear cache entries (pub/sub invoke)

## scribbles

http://nginx.org/en/docs/http/ngx_http_proxy_module.html 

> A cached response is first written to a temporary file, and then the file is renamed. Starting from version 0.8.9, temporary files and the cache can be put on different file systems. However, be aware that in this case a file is copied across two file systems instead of the cheap renaming operation. It is thus recommended that for any given location both cache and a directory holding temporary files are put on the same file system. The directory for temporary files is set based on the use_temp_path parameter (1.7.10). If this parameter is omitted or set to the value on, the directory set by the proxy_temp_path directive for the given location will be used. If the value is set to off, temporary files will be put directly in the cache directory.


If use_temp_path=off, nginx first creates a temporary cache file, e.g.

```
2018/11/26 22:25:14 Got event: notify.Create: "/srv/cache/gw-video-prod/6/1e/8fdbd4af2dbf2d73c6b0b0ed9bbda1e6.0000011462"
2018/11/26 22:25:14 /srv/cache/gw-video-prod/6/1e/8fdbd4af2dbf2d73c6b0b0ed9bbda1e6.0000011462
2018/11/26 22:25:14 &{Wd:2541 Mask:256 Cookie:0 Len:0}
2018/11/26 22:25:14 Got event: notify.Create: "/srv/cache/gw-video-prod/6/1e/8fdbd4af2dbf2d73c6b0b0ed9bbda1e6"
2018/11/26 22:25:14 /srv/cache/gw-video-prod/6/1e/8fdbd4af2dbf2d73c6b0b0ed9bbda1e6
2018/11/26 22:25:14 &{Wd:2541 Mask:128 Cookie:14966 Len:0}
2018/11/26 22:25:14 Got event: notify.Rename: "/srv/cache/gw-video-prod/6/1e/8fdbd4af2dbf2d73c6b0b0ed9bbda1e6.0000011462"
2018/11/26 22:25:14 /srv/cache/gw-video-prod/6/1e/8fdbd4af2dbf2d73c6b0b0ed9bbda1e6.0000011462
2018/11/26 22:25:14 &{Wd:2541 Mask:64 Cookie:14966 Len:0}
```

then renames it to its final destination, better just use IN_MOVED_TO inotify event then:

`man inotify`:
> IN_MOVED_TO (+)
>      Generated for the directory containing the new filename when a file is renamed.

```
2018/11/26 22:41:23 Got event: notify.InMovedTo: "/srv/cache/gw-video-prod/c/c4/7abca09e337e4c00a98266c7efbcbc4c"
2018/11/26 22:41:23 /srv/cache/gw-video-prod/c/c4/7abca09e337e4c00a98266c7efbcbc4c
2018/11/26 22:41:23 &{Wd:833 Mask:128 Cookie:15019 Len:0}
2018/11/26 22:41:27 Got event: notify.InMovedTo: "/srv/cache/gw-video-prod/2/cd/a196e1b9305febb168f2b084f00eccd2"
2018/11/26 22:41:27 /srv/cache/gw-video-prod/2/cd/a196e1b9305febb168f2b084f00eccd2
2018/11/26 22:41:27 &{Wd:3394 Mask:128 Cookie:15020 Len:0}
2018/11/26 22:41:29 Got event: notify.InMovedTo: "/srv/cache/gw-video-prod/a/5d/8db4a84f94cb8d1f40b48159eaec65da"
2018/11/26 22:41:29 /srv/cache/gw-video-prod/a/5d/8db4a84f94cb8d1f40b48159eaec65da
2018/11/26 22:41:29 &{Wd:1450 Mask:128 Cookie:15021 Len:0}
2018/11/26 22:41:32 Got event: notify.InMovedTo: "/srv/cache/gw-video-prod/9/af/75ce461a9954f1c4cbe8dba8d30e3af9"
2018/11/26 22:41:32 /srv/cache/gw-video-prod/9/af/75ce461a9954f1c4cbe8dba8d30e3af9
2018/11/26 22:41:32 &{Wd:1625 Mask:128 Cookie:15022 Len:0}
2018/11/26 22:41:35 Got event: notify.InMovedTo: "/srv/cache/gw-video-prod/b/19/ccdce86106fbc4f1a9b7613b6752119b"
2018/11/26 22:41:35 /srv/cache/gw-video-prod/b/19/ccdce86106fbc4f1a9b7613b6752119b
2018/11/26 22:41:35 &{Wd:1261 Mask:128 Cookie:15023 Len:0}
2018/11/26 22:41:39 Got event: notify.InMovedTo: "/srv/cache/gw-video-prod/8/66/3038a651e89baeb878ddaa220e004668"
2018/11/26 22:41:39 /srv/cache/gw-video-prod/8/66/3038a651e89baeb878ddaa220e004668
2018/11/26 22:41:39 &{Wd:1955 Mask:128 Cookie:15024 Len:0}
```

TODO: 
* check how use_temp_path=on handles writing to cache if temp path and cache dir are not on the same fs