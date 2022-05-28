function loadConfigForm() {
    console.log(serviceName)
    document.getElementById('loadingIcon').style.visibility = "visibile";
    document.getElementById('configForm').style.visibility = "hidden";
    fetch(`/service/${serviceName}/config`).then((res) => {
        if (!res.ok) {
            throw new Error(`HTTP error! Status: ${ res.status }`);
        }
        return res.json()
    }).then(res => {
        console.log(res)
        document.getElementById('configValues').innerHTML = ''
        Object.keys(res).forEach(key => {
            const val = res[key]
            let html = `<label for="${key}" class="form-label">${key}</label>
 <input type="text" class="form-control" name=${key} id="${key}" value="${val}">`
            document.getElementById('configValues').innerHTML += html
        })
        document.getElementById('loadingIcon').style.visibility = "hidden";
        document.getElementById('configForm').style.visibility = "visible";
    })
}

document.addEventListener('submit', (e) => {
    // Store reference to form to make later code easier to read
    const form = e.target
    const data = new FormData(form);
    let json = {};
    data.forEach((value, key) => json[key] = value);
    const body = JSON.stringify(json);

    // Post data using the Fetch API
    fetch(form.action, {
        method: form.method,
        body: body,
    }).then(loadConfigForm);

    // Prevent the default form submit
    e.preventDefault();
});

loadConfigForm()