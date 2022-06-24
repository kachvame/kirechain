<script context='module'>
  // export const load = async ({ fetch }) => {
  // const stats = await fetch('/stats').then((res) => res.json());
  //
  // return {
  //   props: {
  //     topScorers: stats.slice(0, 3),
  //     rows: stats.slice(3).map((stat) => ({
  //       ...stat,
  //       username: stat.username.split('#')[0],
  //       id: stat.username,
  //     })),
  //   },
  // };
  // };
</script>

<script>
  import 'carbon-components-svelte/css/g90.css';
  import Header from '../components/Header.svelte';
</script>

<svelte:head>
  <title>kirechain</title>
</svelte:head>
<Header />
