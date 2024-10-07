const linksListElem = document.getElementById('linksList')

const u = document.getElementById('id').innerText

let allLinks = {}

async function fetchAndDisplayAllLinks() {
    const res = await fetch(`/api/v1/u/${u}/links`) // {"alias1": "url1", "alias2": "url2"}
    if (res.ok) {
        allLinks = await res.json()
        linksListElem.innerHTML = ''
        for (let alias in allLinks) {
            displayOneLinkEnd(alias, allLinks[alias])
        }
    }
}

function displayOneLinkEnd(alias, url) {
    linksListElem.insertAdjacentHTML('beforeend', genLinkDisplayHTML(alias, url))
}
function displayOneLinkBegin(alias, url) {
    linksListElem.insertAdjacentHTML('afterbegin', genLinkDisplayHTML(alias, url))
}

function genLinkDisplayHTML(alias, url) {
    return `<li id="link_${alias}">${location.host}/u/${u}/<span id="alias_${alias}">${alias}</span>
    <button id="editAlias_${alias}" onclick="editAlias('${alias}')">Edit</button>
    <span id="url_${alias}"><a href="${url}">${url}</a></span>
    <button id="editUrl_${alias}" onclick="editUrl('${alias}')">Edit</button>
    <button id="deleteAlias_${alias}" onclick="deleteAlias('${alias}')">Delete</button></li>`
}

function editAlias(alias) {
    document.getElementById(`alias_${alias}`).outerHTML = `<input type="text" value="${alias}" id="aliasEdit_${alias}">`
    document.getElementById(`editAlias_${alias}`).outerHTML = `<button id="editAliasCommit_${alias}" onclick="editAliasCommit('${alias}')">Commit</button> <button id="editAliasCancel_${alias}" onclick="editAliasCancel('${alias}')">Cancel</button>`
}

async function editAliasCommit(oldAlias) {
    const newAlias = document.getElementById(`aliasEdit_${oldAlias}`).value
    const res = await fetch(`/api/v1/u/${u}/${oldAlias}`, {
        method: 'PATCH',
        body: JSON.stringify({
            alias: newAlias
        })
    })
    if (res.ok) {
        updateLink(oldAlias, newAlias, '')
    }
}

function editAliasCancel(oldAlias) {
    updateLink(oldAlias, '', '')
}

function editUrl(alias) {
    document.getElementById(`url_${alias}`).outerHTML = `<input type="text" value="${allLinks[alias]}" id="urlEdit_${alias}">`
    document.getElementById(`editUrl_${alias}`).outerHTML = `<button id="editUrlCommit_${alias}" onclick="editUrlCommit('${alias}')">Commit</button> <button id="editUrlCancel_${alias}" onclick="editUrlCancel('${alias}')">Cancel</button>`
}

async function editUrlCommit(oldAlias) {
    const newUrl = document.getElementById(`urlEdit_${oldAlias}`).value
    const res = await fetch(`/api/v1/u/${u}/${oldAlias}`, {
        method: 'PATCH',
        body: JSON.stringify({
            url: newUrl
        })
    })
    if (res.ok) {
        const serverUrl = await res.text()
        updateLink(oldAlias, '', serverUrl)
    }
}

function editUrlCancel(oldAlias) {
    updateLink(oldAlias, '', '')
}

function deleteAlias(alias) {
    document.getElementById(`deleteAlias_${alias}`).outerHTML = `<button id="deleteAliasCommit_${alias}" onclick="deleteAliasCommit('${alias}')" disabled>Delete: click to confirm</button>`
    setTimeout(() => {
        document.getElementById(`deleteAliasCommit_${alias}`).disabled = false
    }, 1000)
    setTimeout(() => {
        const commitButton = document.getElementById(`deleteAliasCommit_${alias}`)
        if (commitButton !== null) {
            commitButton.outerHTML = `<button id="deleteAlias_${alias}" onclick="deleteAlias('${alias}')">Delete</button>`
        }
    }, 3000)
}

async function deleteAliasCommit(oldAlias) {
    const res = await fetch(`/api/v1/u/${u}/${oldAlias}`, {
        method: 'DELETE'
    })
    if (res.ok) {
        document.getElementById(`link_${oldAlias}`).remove()
    }
}

function updateLink(oldAlias, newAlias, newUrl) {
    if (newAlias === '') {
        newAlias = oldAlias
    } else {
        allLinks[newAlias] = allLinks[oldAlias]
        delete allLinks[oldAlias]
    }
    if (newUrl === '') {
        newUrl = allLinks[newAlias]
    } else {
        allLinks[newAlias] = newUrl
    }
    document.getElementById(`link_${oldAlias}`).outerHTML = genLinkDisplayHTML(newAlias, newUrl)
}

const newLinkAliasElem = document.getElementById('newLinkAlias')
const newLinkUrlElem = document.getElementById('newLinkUrl')

async function addLink() {
    const newAlias = newLinkAliasElem.value
    const newUrl = newLinkUrlElem.value
    const res = await fetch(`/api/v1/u/${u}/${newAlias}`, {
        method: 'POST',
        body: JSON.stringify({
            alias: newAlias,
            url: newUrl
        })
    })
    if (res.ok) {
        const serverUrl = await res.text()
        allLinks[newAlias] = serverUrl
        displayOneLinkBegin(newAlias, serverUrl)
        newLinkAliasElem.value = ''
        newLinkUrlElem.value = ''
    }
}

fetchAndDisplayAllLinks()