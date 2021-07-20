


<form>
    <label for="email">Email</label>
    <input required type="email" id="email" />

    <label for="password">Password</label>
    <input required type="password" id="password" />

    <button type="submit">Create account</button>
</form>




<script>
    async function handleSubmit(event) {
        console.log(event);
        console.log(event.target);
        console.log(event.target.email.value);
        console.log(event.target.password.value);
    }
</script>

