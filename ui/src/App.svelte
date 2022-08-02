<script>
	import Router from 'svelte-spa-router'
	import Home from "./routes/Home.svelte";
	import About from "./routes/About.svelte";
	const routes = {
		// Exact path
		"/": Home,
		"/about": About,
	};
</script>

<main>
	<h1>Yock</h1>
	<Router {routes} />
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
