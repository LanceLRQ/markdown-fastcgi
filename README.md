# Markdown FastCGI Service

## How to use?

1. ⬇️ Clone this repository from github 
2. 📀 Install golang compiler `1.11+`
3. 🔨 Build

    ```go build main.go```
    
4. ⏩ Run it!

    ```./main [-t address]```
   
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
            fastcgi_pass 127.0.0.1:9001
        }
    }
    ``` 
    
🌊 Have fun!