<script lang="ts">
  import {ListNodes, GetCurrentContext} from '../wailsjs/go/main/App.js'
  import type {k8s} from '../wailsjs/go/models'
  import {BrowserOpenURL, ClipboardSetText} from "../wailsjs/runtime";
  import TailwindCSS from "./lib/TailwindCSS.svelte";

  let currentContext: string = "Unknown"
  let loading: boolean = false
  let nodes: k8s.Node[] = []
  let message: string = ""

  function listNodes(): void {
    loading = true
    ListNodes().then(result => {
      nodes = result
      loading = false
    })
  }

  function getCurrentContext(): void {
    GetCurrentContext().then(result => {
      currentContext = result
    })
  }

  function copyNodeName(name: string): void {
    ClipboardSetText(name).then(_ => {
      message = `Node name copied to clipboard: ${name}`
      setTimeout(() => message = "", 5000)
    })
  }

  getCurrentContext()
  listNodes()
</script>

<header>
  <span>Current Context:</span> {currentContext}
</header>

<main>
  {#if loading}
  <p>Loading nodes ...</p>
  {/if}

  {#if message}
  <div class="absolute flex bottom-0 items-center bg-blue-300 text-white text-sm font-bold px-4 py-3" role="alert">
    <svg class="fill-current w-4 h-4 mr-2" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M12.432 0c1.34 0 2.01.912 2.01 1.957 0 1.305-1.164 2.512-2.679 2.512-1.269 0-2.009-.75-1.974-1.99C9.789 1.436 10.67 0 12.432 0zM8.309 20c-1.058 0-1.833-.652-1.093-3.524l1.214-5.092c.211-.814.246-1.141 0-1.141-.317 0-1.689.562-2.502 1.117l-.528-.88c2.572-2.186 5.531-3.467 6.801-3.467 1.057 0 1.233 1.273.705 3.23l-1.391 5.352c-.246.945-.141 1.271.106 1.271.317 0 1.357-.392 2.379-1.207l.6.814C12.098 19.02 9.365 20 8.309 20z"/></svg>
    <p>{message}</p>
  </div>
  {/if}

  {#if nodes}
  <ul id="nodes" class="mt-4">
    {#each nodes as {name, consolePageURL, dashboardURL}}
      <li class="group mx-4 my-4 flex rounded-2xl border px-4 py-4 hover:bg-cyan-50 justify-between">
        <div class="flex-initial w-30" on:click={copyNodeName(name)}>{name}</div>
        <div class="flex justify-between hidden group-hover:block">
          <a on:click={BrowserOpenURL(consolePageURL)} class="hover:cursor-pointer">AWS EC2 Page</a>
          <a on:click={BrowserOpenURL(dashboardURL)} class="hover:cursor-pointer">Datadog Host Dashboard</a>
        </div>
      </li>
    {/each}
  </ul>
  {/if}
</main>

<style>
</style>
