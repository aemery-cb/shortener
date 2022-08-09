<script>
	import Router from "svelte-spa-router";
	import Home from "./routes/Home.svelte";
	import About from "./routes/About.svelte";

	import { isAuthenticated, user, token } from "./store";
	import auth from "./authService";
	import { onMount } from "svelte";

	const routes = {
		// Exact path
		"/": Home,
		"/about": About,
	};

	let auth0Client
	onMount(async () => {
		auth0Client = await auth.createClient();

		isAuthenticated.set(await auth0Client.isAuthenticated());
		user.set(await auth0Client.getUser());
	});

	async function login() {
		await auth.loginWithPopup(auth0Client);
		token.set(await auth0Client.getTokenSilently())
	}
	function logout() {
		return auth0Client.logout({
			returnTo: window.location.hostname 
		});
	}

</script>

<main>
	<h1>Yock</h1>
	<Router {routes} />
	{#if $isAuthenticated}
	<button on:click={logout}>Logout</button>
	{:else}
	<button on:click={login}>Login</button>
	{/if}
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
