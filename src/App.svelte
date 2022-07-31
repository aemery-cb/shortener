<script>
	let newUrl = "";
	async function submit (e) {

		const formData = new FormData(e.target);

		const data = {};
		for (let field of formData) {
			const [key, value] = field;
			data[key] = value;
		}
		console.log(data);
		newUrl = await shortenURL(data.url);
	}

	const shortenURL = async (url) => {
		const response = await fetch('/api/shorten', {
			method: "POST",
			body: JSON.stringify({
				url: url
			})
		})

		const data = await response.json()
		return data.url;
	}
</script>

<main>
	<form on:submit|preventDefault={submit}>
		<input type="text" id="url" name="url" />
		<button type="submit"> Shorten! </button>
	</form>

	{#if newUrl != ""}
	<div>{newUrl}</div>
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
