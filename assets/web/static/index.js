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
    }).then(() =>
        window.location.reload()
    );

    // Prevent the default form submit
    e.preventDefault();
});

console.log("hey!");

function clearScreen() {
    fetch("/screens/current", {
        method: "DELETE"
    }).then(() => window.location.href=window.location.href).catch(err => console.log(err));
}
