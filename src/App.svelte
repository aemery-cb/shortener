<script>
	let newUrl = "";
    function validURL(str) {
  var pattern = new RegExp('^(https?:\\/\\/)?'+ // protocol
    '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|'+ // domain name
    '((\\d{1,3}\\.){3}\\d{1,3}))'+ // OR ip (v4) address
    '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*'+ // port and path
    '(\\?[;&a-z\\d%_.~+=-]*)?'+ // query string
    '(\\#[-a-z\\d_]*)?$','i'); // fragment locator
  return !!pattern.test(str);
}

	async function submit (e) {

		const formData = new FormData(e.target);

		const data = {};
		for (let field of formData) {
			const [key, value] = field;
			data[key] = value;
		}

        if (validURL(data.url)) {
		    newUrl = await shortenURL(data.url);
        } else {
            newUrl = "please enter a valid url"
        }
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
<h1>Yock</h1>
<p>Enter a url to shorten below</p>
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
