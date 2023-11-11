<script lang="ts">
  import {ListNodes, GetCurrentContext} from '../wailsjs/go/main/App.js'
  import type {k8s} from '../wailsjs/go/models'
  import {BrowserOpenURL} from "../wailsjs/runtime";
  import TailwindCSS from "./lib/TailwindCSS.svelte";

  let currentContext: string = "Unknown"
  let loading: boolean = false
  let nodes: k8s.Node[] = []

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

  {#if nodes}
  <ul id="nodes" class="mt-4">
    {#each nodes as {name, consolePageURL, dashboardURL}}
      <li class="group mx-4 my-4 flex rounded-2xl border px-4 py-4 hover:bg-cyan-50 justify-between">
        <div class="flex-initial w-30">{name}</div>
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
