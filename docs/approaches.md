# Approaches to dealing with CLI flags

## Package librarian

A search on [`pkg.go.dev`](https://pkg.go.dev/search?q=flag&m=) provides some insight about the many CLI flag parsing
packages out there.

| repository                             |:Imported by`*`:| Approach     | Distinctive features                              | Last published |
-----------------------------------------------------------------------------------------------------------------------------------------------
|`flag` (standard library)               | 278,354        | programmatic | Dead simple. Supports '-flag' or '--flag' styles. | 2023/09/06 |
|`https://github.com/spf13/pflag`        | 29,126         | programmatic | Replaces `flag`. Supports POSIX style, e.g. '-xvf --flag'. | 2019/09/18 |
|                                        |                |              | Arrays, maps & extra types. Extensible.           |            |
|`https://github.com/jessevdk/go-flags`  | 8,230          | struct tags  | Part of a full tag-driven command builder.        | 2021/03/21 |
|                                        |                |              | Rich set of tag decorators.                       |            |
|`https://github.com/ogier/pflag`        | 738            | programmatic | Original repo spf13/pflag was forked from         | 2015/02/15 |
|`https://github.com/mreiferson/go-options` | 325         | struct tags  | Tag-driven CLI flag parser                        | 2019/12/25 |
|`https://github.com/juju/gnuflag`       | 246            | programmatic | Replaces `flag`. Supports GNU style.              | 2017/11/13 |
|`https://github.com/projectdiscovery/goflags` | 183      | programmatic | Extend `flag` with goodies                        | 2023/09/18 |
|`https://github.com/rclone/rclone/fs/config/flags`| 101  | programmatic | Extend `spf13/pflag` to support env. variables.   | 2023/09/11 |
|`https://github.com/ViBiOh/flags`       | 57             | programmatic | Zero-dependency flag parser with env. variables.  | 2023/08/26 |
|`https://github.com/btcsuite/go-flags`  | 17             | struct tags  | Tag-driven CLI flag parser                        | 2015/01/16 |
|`https://github.com/integrii/flaggy`    | 11             | programmatic | Part of a simple command builder.                 | 2022/05/28 |
|`https://github.com/podhmo/flagstruct`  | 3              | struct tags  | Brings tags support to `flag`, `pflag`            | 2022/08/14 |
|`https://github.com/fredbi/gflag`       | 1              | programmatic | Brings generics to `pflag`                        | 2023/09/26 |
|`https://github.com/AdamSLevy/flagbind` | 0              | struct tags  | Brings tags support to `flag`, `pflag`            | 2020/07/35 |


`*`: reported by `pkg.go.dev` as of 2023/09/26

## Comments

> Comments by fredbi, as of 2023/09/26

If we exclude the `pflag` standard library, usage of which is boosted by the eco-system of standard libraries,
`spf13/pflag` seems the most popular package,
most likely due to the popularity of the companion packages `spf13/cobra` and `spf13/viper`.

The programmatic approach seems overall favored by go developers.

I've personally used `github.com/jessevdk/go-flags` (N.B: as part of the [`go-swagger`](https://github.com/go-swagger/go-swagger)
CLI): it is really great for simple CLIs, for which a low-code approach is desirable.

When it comes to maintain projects with a many commands, subcommands, flags etc. over a long period of time (this includes `go-swagger`),
we're all tempted to move to a more explicit, no-magic kind of approach.

Further, parsing flags is certainly not the end of the road: what is likable about `spf13/viper` for instance, is how
it transforms your app into a full-fledged 12-factor app with no effort.

Now `spf13/pflag` certainly has some shortcomings. Most importantly, it has become poorly maintained (ok, it just works, but still)
and clearly lacks modern improvements (for instance, like what `gflag` is suggesting).

An interesting finding from the above table is the recent development of "zero-dependency" clones with support for 
environment variables (well, `spf13/pflag` doesn't add dependencies either). I assume that this in reaction to the many
dependencies imported with `spf3/viper`. I've never found dependencies much of a problem myself, at least since go modules 
are out.
