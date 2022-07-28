import logo from './logo.svg';
import styles from './App.module.css';
import { lazy, createResource, onCleanup } from "solid-js";

import ColorGrid from './ColorGrid'

const fetchGrid = async () =>
  (await fetch('http://localhost:3030/')).json()

function App() {
  const [grid, {_, refetch}] = createResource(fetchGrid)

  const refresher = setInterval(refetch, 500)
  onCleanup(() => clearInterval(refresher));

  return (
    <>
      <div class="box-container">
        <For each={ grid() }>{(box, i) => <div style={`background-color: ${box};`} class={"box " + box}></div>}
        </For>
      </div>
    </>
  );
}

export default App;
