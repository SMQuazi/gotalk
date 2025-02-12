import "./App.css";
import { Signal, useSignal } from "use-signals";

const messager = new Signal.State("");
const user = new Signal.State<string | undefined>(undefined);

let socket: WebSocket;

function App() {
  const message = useSignal(messager);
  const userInfo = useSignal(user);

  const handleSubmit = (e) => {
    e.preventDefault();
    e.stopPropagation();
    console.log(e.target[0].value);
    user.set(e.target[0].value);

    socket = new WebSocket(
      `ws://localhost:8080/ws?authorization=${user.get()}`
    );
    socket.addEventListener("close", () => {
      console.log("Socket closed");
    });
    socket.addEventListener("error", (error) => {
      console.error("Socket error", error);
    });
    socket.addEventListener("message", (event) => {
      messager.set(event.data);
    });
  };

  const handleChange = (e) => {
    console.log(e.target.value);
    socket.send(e.target.value);
  };

  return (
    <>
      {userInfo ? (
        <div>
          <h2>Message: {message}</h2>
          <input type="text" onChange={handleChange} />
        </div>
      ) : (
        <form onSubmit={handleSubmit}>
          <div>
            <h2>Enter your name</h2>
            <input type="text" onBlur={(e) => user.set(e.target.value)} />
          </div>
        </form>
      )}
    </>
  );
}

export default App;
