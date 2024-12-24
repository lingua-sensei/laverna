# Laverna

the goddess of the thieves, helps you steal translation speeches from the G-daddy monopoly.

<img src="https://github.com/user-attachments/assets/d1d344c9-f36b-4cf7-af70-f162f93ea9f0" width="400" alt="squash-goose">

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
laverna -file example.yaml 
```


> [!IMPORTANT]  
> There seems to be heavy obfuscatication over 200 characters, I am working on an algo that to crack that.
