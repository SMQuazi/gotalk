import "./App.css";
import { Signal, useSignal } from "use-signals";

const users = new Signal.State({});
const thisUser = `user${Math.floor(Math.random() * 90 + 10)}`;

const socket = new WebSocket(
  `ws://localhost:8080/ws?authorization=${thisUser}`
);
socket.addEventListener("close", () => {
  console.log("Socket closed");
});
socket.addEventListener("error", (error) => {
  console.error("Socket error", error);
});
socket.addEventListener("message", (event) => {
  const json = JSON.parse(event.data);
  users.set(json.users);
});

function App() {
  const userInfo = useSignal(users);

  return (
    <>
      <h1>Welcome {thisUser}</h1>
      <h2>All Users:</h2>
      {Object.keys(userInfo).map((user) => (
        <div key={user}>{user}</div>
      ))}
    </>
  );
}

export default App;
