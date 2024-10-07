const passwordElem = document.getElementById('password')
const passwordConfirmElem = document.getElementById('passwordConfirm')
const confirmMessageElem = document.getElementById('confirmMessage')
const submitElem = document.getElementById('submit')

function confirmPasswordsAreSame(ev) {
    if (passwordElem.value !== passwordConfirmElem.value) {
        console.log(`passwords don't match: ${passwordElem.value} // ${passwordConfirmElem.value}`)
        submitElem.disabled = true
        confirmMessageElem.innerHTML = 'The passwords must match.'
    } else {
        console.log(`passwords match: ${passwordElem.value} // ${passwordConfirmElem.value}`)
        submitElem.disabled = false
        confirmMessageElem.innerHTML = ''
    }
}

passwordElem.addEventListener('change', confirmPasswordsAreSame)
passwordConfirmElem.addEventListener('change', confirmPasswordsAreSame)

const formElem = document.getElementById('signupForm')

async function signUp() {
    const formData = new FormData(formElem)
    const bodyJSON = {
        id: formData.get('id'),
        email: formData.get('email'),
        password: formData.get('password')
    }
    console.log(bodyJSON)
    const res = await fetch('/signup/post', {
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
    signUp()
})