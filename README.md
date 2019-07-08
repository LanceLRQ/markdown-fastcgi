# Markdown FastCGI Service

Parse markdown document to html in FastCGI supported web server like `nginx`.

## 🤔 How to use?

1. ⬇️ Clone this repository from github.
2. 📀 Install golang compiler `1.11+`.
3. 🔨 Build it!

    ```go build main.go```
    
4. ⏩ Run it!

    ```./main [-l address]```
   
   Example:
    ```
    ./main
    ./main -l 127.0.0.1:9001
    ./main --listen 127.0.0.1:9002
    ./main -l unix://var/run/mds.sock
    ```
    If your command arguments miss `-l`，it will listen at `127.0.0.1:9001` 
    
    You can use `supervisor` to hold this service.
    
5. 🛠 Configure your nginx.
    
    ```
    server {
        # your configurations.
        
        location ~ \.md$ {
            include fastcgi_params;
            fastcgi_pass 127.0.0.1:9001;
        }
    }
    ``` 

6. 🌊 Visit your website.

## 🙏 Thanks

[Blackfriday: a markdown processor for Go](github.com/russross/blackfriday)

### 😜 Have fun!🌞