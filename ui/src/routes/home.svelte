<script>
  import {link} from 'svelte-spa-router'
	let newUrl = "";
	function validURL(str) {
		var pattern = new RegExp(
			"^(https?:\\/\\/)?" + // protocol
				"((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|" + // domain name
				"((\\d{1,3}\\.){3}\\d{1,3}))" + // OR ip (v4) address
				"(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*" + // port and path
				"(\\?[;&a-z\\d%_.~+=-]*)?" + // query string
				"(\\#[-a-z\\d_]*)?$",
			"i"
		); // fragment locator
		return !!pattern.test(str);
	}

	async function submit(e) {
		const formData = new FormData(e.target);
		const data = {};
		for (let field of formData) {
			const [key, value] = field;
			data[key] = value;
		}

		if (validURL(data.url)) {
			newUrl = await shortenURL(data.url);
		} else {
			newUrl = "please enter a valid url";
		}
	}

	const shortenURL = async (url) => {
		const response = await fetch("/api/shorten", {
			method: "POST",
			headers: {
				"content-type": "application/json"
			},
			body: JSON.stringify({
				url: url,
			}),
		});

		const data = await response.json();
		return data.url;
	};
</script>

<main>
	<p>Enter a url to shorten below</p>
	<form on:submit|preventDefault={submit}>
		<input type="text" id="url" name="url" />
		<button type="submit"> Shorten! </button>
	</form>

	{#if newUrl != ""}
		<div>{newUrl}</div>
	{/if}

  <a href="/about" use:link>About</a>
</main>

