# Messenger

Using on your local.

After `go run`

Using with `python` `requests`:

```
In [7]: import requests 

In [8]: temp = requests.put("http://localhost:8090/send_message", data='{"BobToJohn": "Hello"}')

In [9]: temp.content
Out[9]: ''

In [10]: temp = requests.get("http://localhost:8090/get_message", params={"message_key": "BobToJohn"})

In [11]: temp.content
Out[11]: 'Hello'

```

Using with linux shell

```
curl -X POST -d '{"BobToJohn": "Hello"}' http://localhost:8090/send_message
message=$(curl -X GET http://localhost:8090/get_message?message_key=BobToJohn)
echo message
```

