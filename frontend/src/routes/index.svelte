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

<style lang='scss'>
  @import '@carbon/colors';

  .container {
    display: flex;
    justify-content: center;
    padding: 1rem;
    margin-left: auto;
    margin-right: auto;
  }

  @media (max-width: 420px) {
    .container {
      flex-direction: column;
    }

    :global(.bx--btn) {
      max-width: none;
    }
  }


  h1 {
    text-align: center;
    overflow-wrap: break-word;
  }

  .divider {
    height: 2px;
    width: 100%;
    background: $gray-80;
    margin-top: 1rem;
    margin-bottom: 1rem;
  }

  .page {
    max-width: 580px;
    margin-left: auto;
    margin-right: auto;
  }
</style>

<svelte:head>
  <title>kirechain</title>
</svelte:head>
<Header />
<div class='page'>
  <div class='container'>
    <TextInput placeholder='Start of message' size='xl' bind:value={input} />
    <Button iconDescription='Go' on:click={handleClick} icon={NextOutline} href='#' kind='secondary'>
      Generate
    </Button>
  </div>
  {#each messages as message}
    <h1>
      {message}
    </h1>
    <div class='divider' />
  {/each}
</div>
