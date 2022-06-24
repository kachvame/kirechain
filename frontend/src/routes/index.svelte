<script context='module'>
  export const load = async ({ fetch }) => {
    const response = await fetch('/chain').then((res) => res.text());

    return {
      props: {
        messages: [response],
      },
    };
  };
</script>

<script>
  import 'carbon-components-svelte/css/g90.css';
  import Header from '../components/Header.svelte';
  import { TextInput } from 'carbon-components-svelte';
  import { Button } from 'carbon-components-svelte';
  import NextOutline from 'carbon-icons-svelte/lib/NextOutline.svelte';

  export let messages;
  let input;


  async function handleClick() {
    const params = new URLSearchParams();
    if (input !== '') {
      params.set('input', input);
    }
    const response = await fetch(`/chain?${params.toString()}`).then(v => v.text());
    messages = [response, ...messages].slice(0, 10);
  }
</script>
<style>
    .container {
        display: flex;
        justify-content: center;
        padding: 1rem;
        margin-left: auto;
        margin-right: auto;
        width: 512px;
    }

    h1 {
        text-align: center;
    }
</style>

<svelte:head>
  <title>kirechain</title>
</svelte:head>
<Header />
<div class='container'>
  <TextInput placeholder='Start of message' size='xl' bind:value={input} />
  <Button iconDescription='Go' on:click={handleClick} icon={NextOutline} href='#' kind='secondary'>Generate</Button>
</div>
{#each messages as message}
  <h1>
    {message}
  </h1>
{/each}
