function loadDashboard() {
    document.getElementById('currentServiceLoading').style.visibility = "visible";
    document.getElementById('noCurrentService').style.visibility = "hidden";
    document.getElementById('currentService').style.visibility = "hidden";
    loadCurrentService()
    // fetch("/services").then((res) => {
    //     if (!res.ok) {
    //         throw new Error(`HTTP error! Status: ${ res.status }`);
    //     }
    //     return res.json()
    // }).then((res) => {
    //     console.log(`Received response from services endpoint: ${JSON.stringify(res)}`)
    //     if (!res.success) {
    //         console.error("Received unsuccessful response from services endpoint")
    //     }
    //     let services = res.response
    //     services.forEach(elem => {
    //         let html = `<option value=${elem}>${elem}</option>`
    //         document.getElementById('serviceSelector').innerHTML += html
    //     })
    //     document.getElementById('loadingIcon').style.visibility = "hidden";
    //     document.getElementById('serviceForm').style.visibility = "visible";
    // })
}


function loadCurrentService() {
    document.getElementById('currentServiceLoading').style.visibility = "visible";
    fetch('/service', {
        method: "GET"
    }).then(res => {
        if (!res.ok) {
            throw new Error(`HTTP error! Status: ${ res.status }`);
        }
        return res.json()
    }).then(json => {
        if (!json.success) {
            throw new Error('API operation was unsuccessful')
        }
        document.getElementById('currentServiceLoading').style.visibility = "hidden";
        if (json.response['is_running']) {
            const svc = document.getElementById('currentService')
            svc.innerHTML = json.response['id']
            svc.style.visibility = "visible";
        } else {
            document.getElementById('no').style.visibility = "hidden";
        }
    }).catch(err => {
        document.getElementById('currentServiceLoading').style.visibility = "hidden";
        const serviceError = document.getElementById('serviceError')
        serviceError.innerHTML = `Unable to fetch: ${err}`
        serviceError.style.visibility = "visible";
    })
}


function stopService() {
    fetch("/service", {
        method: 'DELETE',
    }).then(loadServices);
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

loadDashboard()