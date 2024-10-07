const formElem = document.getElementById('loginForm')

async function logIn() {
    const formData = new FormData(formElem)
    const bodyJSON = {
        id: formData.get('id'),
        password: formData.get('password')
    }
    console.log(bodyJSON)
    const res = await fetch('/login/post', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(bodyJSON)
    })
    if (res.ok) {
        const json = await res.json()
        console.log(json)
        const jwtCookie = json['jwtCookie']
        console.log(`jwtCookie: ${jwtCookie}`)
        document.cookie = jwtCookie
        window.open('/myaccount', '_self')
    }
}

formElem.addEventListener('submit', (event) => {
    event.preventDefault()
    logIn()
})