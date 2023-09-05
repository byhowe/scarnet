# scarnet

Scarnet is a messaging application written in Golang. The name is an omage to a
server I wrote back in 2020 called Scarlett. I am currently learning Go and I
thought it would be a good excercise to rewrite the server in Go. I had tried
not to use too many libraries for Scarlett and I am trying to do the same with
Scarnet. I have not even used http. The communication happens through TCP
sockets that the messages pass through after being serialized using JSON. I have
no plans to write a GUI for a client, but I do plan to create a client library.

---

## License

The project is licensed under MIT.

---

If you have any ideas or recommendations, please feel free to share. Afterall,
we are all here for learning.
