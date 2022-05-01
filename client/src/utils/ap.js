import logError from './errors';
import { proxyClient } from './http';

const apHeaders = {
    // accept: 'application/ld+json; profile="https://www.w3.org/ns/activitystreams"',
    accept: 'application/activity+json',
    // contentType: 'application/ld+json; profile="https://www.w3.org/ns/activitystreams"'
    contentType: 'application/activity+json'
}

var activityTypes = ["Accept", "Add", "Announce", "Arrive", "Block", "Create", "Delete", "Dislike", "Flag", "Follow", "Ignore", "Invite", "Join", "Leave", "Like", "Listen", "Move", "Offer", "Question", "Reject", "Read", "Remove", "TentativeReject", "TentativeAccept", "Travel", "Undo", "Update", "View"]
var actorTypes = ["Application", "Group", "Organization", "Person", "Service"]
var objectTypes = ["Article", "Audio", "Document", "Event", "Image", "Note", "Page", "Place", "Profile", "Relationship", "Tombstone", "Video"]
var linkTypes = ["Mention"]
// var audiences = ["to", "bto", "cc", "bcc", "audience"]

function isValidActivity(activity) {
    return activity.id && activityTypes.includes(activity.type);
}

function isValidActor(actor) {
    return actor.id && actorTypes.includes(actor.type);
}

function isValidObject(object) {
    return object.id && objectTypes.includes(object.type);
}

function isValidLink(link) {
    return link.id && linkTypes.includes(link.type);
}

async function fetchCollectionWithItems(iri, withCredentials=false) {
    const collection = await fetchCollection(iri, withCredentials);
    if (!collection.orderedItems) {
        collection.orderedItems = [];
    }
    if (collection.first) {
        const first = typeof collection.first !== 'string' ? collection.first.id : collection.first;
        let page = await fetchPage(first, withCredentials);
        if (page && (page.orderedItems || page.items)) {
            collection.orderedItems = collection.orderedItems.concat(page.orderedItems || page.items);
        }
        let emptyPage = false;
        while (page && (page.next && !emptyPage)) {
            const next = typeof page.next !== 'string' ? page.next.id : page.next;
            page = await fetchPage(next, withCredentials);
            if (page.orderedItems || page.items) {
                collection.orderedItems = collection.orderedItems.concat(page.orderedItems || page.items);
                emptyPage = (page.orderedItems || page.items).length === 0;
            }
        }
    }
    return collection;
}

function fetchCollection(iri, withCredentials=false) {
    return proxyClient.get(iri, {
        headers: {'Accept': apHeaders.accept},
        withCredentials
    } ).then((res) => {
        console.log(res);
        return res.data;
    }).catch(error => {
        logError(error);
    });
}

function fetchPage(iri, withCredentials=false) {
    return proxyClient.get(iri, {
        headers: {'Accept': apHeaders.accept},
        withCredentials
    } ).then(res => {
        console.log(res);
        return res.data;
    }).catch(error => {
        logError(error);
    });
}

function renderIcon(icon) {
    if (!icon) {
        return <img className="Actor-image" src={"/logo192.png"} alt="..." />;
    }
    switch (icon.type) {
        case "Image":
            return <img className="Actor-image" src={icon.url} alt="..." />;
        default:
            break;
    }
}

export {
    apHeaders,
    isValidActivity,
    isValidActor,
    isValidObject,
    isValidLink,
    fetchCollectionWithItems,
    fetchCollection,
    fetchPage,
    renderIcon
}