# Laverna

the goddess of the thieves, helps you steal translation speeches from the G-daddy monopoly.

```
        w*W*W*W*w
         \"."."/
          //`\\
         (/o o\)
         (\_-_/) 
        .-~'='~-.
       /`~`"Y"`~`\
      / /(_ * _)\ \       
     / /  )   (  \ \
     \ \_/\\_//\_/ / 
      \/_) '*' (_\/
        |   L   |
        |   a   |
        |   v   |
        |   e   |
        |   r   |
        |   n   |
        |   a   |
        |       |
        w*W*W*W*w
```

### Install Via Go

```shell
go install github.com/lingua-sensei/laverna@latest
```

### Grab Binaries

You can find binaries through GitHub releases.

### Sample Usage

Let's create example YAML.

```yaml
- speed: normal
  voice: th
  text: "สวัสดีครับ"
- speed: slower
  voice: en
  text: "Hello there"
- speed: slowest
  voice: ja
  text: "こんにちは~"
```

Running below command will generate audios in the same directory.

```shell
go run main.go -file example.yaml 
```


> [!IMPORTANT]  
> There seems to be heavy obfuscatication over 200 characters, I am working on an algo that to crack that.