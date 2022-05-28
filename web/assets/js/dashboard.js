function loadServiceForm() {
    document.getElementById('loadingIcon').style.visibility = "visible";
    document.getElementById('serviceForm').style.visibility = "hidden";
    fetch("/services").then((res) => {
        if (!res.ok) {
            throw new Error(`HTTP error! Status: ${ response.status }`);
        }
        return res.json()
    }).then((res) => {
        console.log(res)
        let selected = res.selected
        let serviceRunning = res.selected !== null
        const link = document.getElementById('serviceConfigLink')
        if (serviceRunning) {
            link.href = "/config/" + res.selected
            link.innerText = selected + " config"
            link.style.visibility = "visible"
            document.getElementById('serviceSelector').innerHTML = ""
        } else {
            document.getElementById('serviceSelector').innerHTML = "<option disabled=\"disabled\" selected></option>"
            link.style.visibility = "hidden"
        }

        res.services.forEach(elem => {
            let html = `<option value=${elem} ${elem === selected ? 'selected' : ''}>${elem}</option>`
            document.getElementById('serviceSelector').innerHTML += html
            console.log(elem)
        })
        document.getElementById('loadingIcon').style.visibility = "hidden";
        document.getElementById('serviceForm').style.visibility = "visible";
    })
}

function stopService() {
    fetch("/service", {
        method: 'DELETE',
    }).then(loadServiceForm);
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
    }).then(loadServiceForm);

    // Prevent the default form submit
    e.preventDefault();
});

loadServiceForm()