<script lang="ts">
  import {ListNodes, GetCurrentContext} from '../wailsjs/go/main/App.js'
  import type {k8s} from '../wailsjs/go/models'
  import {BrowserOpenURL, ClipboardSetText} from "../wailsjs/runtime";
  import TailwindCSS from "./lib/TailwindCSS.svelte";

  let currentContext: string = "Unknown"
  let loading: boolean = false
  let nodes: k8s.Node[] = []
  let message: string = ""

  function listNodes(shouldClearCache: boolean = false): void {
    nodes = []
    loading = true
    ListNodes(shouldClearCache).then(result => {
      nodes = result
      sort("instanceType", false)
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

  let sortBy = {col: "instanceType", ascending: true};

  $: sort = (column: string, shouldInvertAscending: boolean = true) => {

    if (sortBy.col == column && shouldInvertAscending) {
      sortBy.ascending = !sortBy.ascending
    } else {
      sortBy.col = column
      sortBy.ascending = true
    }

    // Modifier to sorting function for ascending or descending
    let sortModifier = (sortBy.ascending) ? 1 : -1;

    let sort = (a, b) =>
      (a[column] < b[column])
      ? -1 * sortModifier
      : (a[column] > b[column])
      ? 1 * sortModifier
      : 0;

    nodes = nodes.sort(sort);
  }

  function humanizeSize(bytes: number): string {
    let quantity: number
    let unit: string
    if (bytes < 2 ** 10) {
      quantity = bytes
      unit = "B"
    } else if (bytes < 2 ** 20) {
      quantity = bytes / 2 ** 10
      unit = "KiB"
    } else if (bytes < 2 ** 30) {
      quantity = bytes / 2 ** 20
      unit = "MiB"
    } else if (bytes < 2 ** 40) {
      quantity = bytes / 2 ** 30
      unit = "GiB"
    } else {
      quantity = bytes / 2 ** 40
      unit = "TiB"
    }
    return `${quantity.toFixed(2)} ${unit}`
  }

  getCurrentContext()
  listNodes()
</script>

<header>
  <div>
    <span class="font-bold">Current Context:</span> {currentContext}
  </div>
  <button
    class="border rounded-2xl px-4 py-2"
    on:click={() => listNodes(true)} >Refresh</button>
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

  {#if nodes.length > 0}
  <table class="mt-4">
    <thead>
      <tr>
        <th on:click={sort("name")}>Node</th>
        <th on:click={sort("instanceType")}>Instance Type</th>
        <th>CPU Usage</th>
        <th>Mem Usage</th>
        <th>EC2</th>
        <th>Datadog</th>
      </tr>
    </thead>
    <tbody>
    {#each nodes as {name, consolePageURL, dashboardURL, instanceType, usage, capacity}}
      <tr class="mx-4 my-4 px-4 py-4 hover:bg-cyan-50 node">
        <td class="name" on:click={copyNodeName(name)}>{name}</td>
        <td>
          <a on:click={BrowserOpenURL(`https://instances.vantage.sh/aws/ec2/${instanceType}`)}>{instanceType}</a>
        </td>
        <td>
          {usage.cpu.toFixed(1)}/{capacity.cpu}
        </td>
        <td>
          {humanizeSize(usage.memory)}/{humanizeSize(capacity.memory)}
        </td>
        <td class="ec2">
          <a on:click={BrowserOpenURL(consolePageURL)}></a>
        </td>
        <td class="datadog">
          <a on:click={BrowserOpenURL(dashboardURL)}></a>
        </td>
      </tr>
    {/each}
    </tbody>
  </table>
  {/if}
</main>

<style>
  .node {
    vertical-align: middle;
  }
  .node a {
    display: inline-block;
    margin-top: 1em;
    margin-left: 1.2em;
    margin-right: 1.2em;
    background-size: cover;
    background-repeat: no-repeat;
  }
  .node a:hover {
    cursor: pointer;
  }

  .node .name {
    max-width: 16em;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .ec2 a {
    width: 36px;
    height: 36px;
    background-image: url('data:image/svg+xml;base64,PHN2ZyBjbGFzcz0idy02IGgtNiIgaGVpZ2h0PSI0MCIgd2lkdGg9IjQwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxkZWZzPjxsaW5lYXJHcmFkaWVudCB4MT0iMCUiIHkxPSIxMDAlIiB4Mj0iMTAwJSIgeTI9IjAlIiBpZD0iQXJjaF9BbWF6b24tRUMyXzMyX3N2Z19fYSI+PHN0b3Agc3RvcC1jb2xvcj0iI0M4NTExQiIgb2Zmc2V0PSIwJSI+PC9zdG9wPjxzdG9wIHN0b3AtY29sb3I9IiNGOTAiIG9mZnNldD0iMTAwJSI+PC9zdG9wPjwvbGluZWFyR3JhZGllbnQ+PC9kZWZzPjxnIGZpbGw9Im5vbmUiIGZpbGwtcnVsZT0iZXZlbm9kZCI+PHBhdGggZD0iTTAgMGg0MHY0MEgweiIgZmlsbD0idXJsKCNBcmNoX0FtYXpvbi1FQzJfMzJfc3ZnX19hKSI+PC9wYXRoPjxwYXRoIGQ9Ik0yNi4wNTIgMjdMMjYgMTMuOTQ4IDEzIDE0djEzLjA1MkwyNi4wNTIgMjd6TTI3IDE0aDJ2MWgtMnYyaDJ2MWgtMnYyaDJ2MWgtMnYyaDJ2MWgtMnYyaDJ2MWgtMnYuMDUyYS45NS45NSAwIDAxLS45NDguOTQ4SDI2djJoLTF2LTJoLTJ2MmgtMXYtMmgtMnYyaC0xdi0yaC0ydjJoLTF2LTJoLTJ2MmgtMXYtMmgtLjA1MmEuOTUuOTUgMCAwMS0uOTQ4LS45NDhWMjdoLTJ2LTFoMnYtMmgtMnYtMWgydi0yaC0ydi0xaDJ2LTJoLTJ2LTFoMnYtMmgtMnYtMWgydi0uMDUyYS45NS45NSAwIDAxLjk0OC0uOTQ4SDEzdi0yaDF2Mmgydi0yaDF2Mmgydi0yaDF2Mmgydi0yaDF2Mmgydi0yaDF2MmguMDUyYS45NS45NSAwIDAxLjk0OC45NDhWMTR6bS02IDE5SDdWMTloMnYtMUg3LjA2MkM2LjQ3NyAxOCA2IDE4LjQ3NyA2IDE5LjA2MnYxMy44NzZDNiAzMy41MjMgNi40NzcgMzQgNy4wNjIgMzRoMTMuODc3Yy41ODUgMCAxLjA2MS0uNDc3IDEuMDYxLTEuMDYyVjMxaC0xdjJ6TTM0IDcuMDYydjEzLjg3NmMwIC41ODUtLjQ3NiAxLjA2Mi0xLjA2MSAxLjA2MkgzMHYtMWgzVjdIMTl2M2gtMVY3LjA2MkMxOCA2LjQ3NyAxOC40NzcgNiAxOS4wNjIgNmgxMy44NzdDMzMuNTI0IDYgMzQgNi40NzcgMzQgNy4wNjJ6IiBmaWxsPSIjRkZGIj48L3BhdGg+PC9nPjwvc3ZnPg==');
  }

  .datadog a {
    width: 36px;
    height: 36px;
    background-image: url('https://www.vectorlogo.zone/logos/datadoghq/datadoghq-icon.svg');
  }
</style>
