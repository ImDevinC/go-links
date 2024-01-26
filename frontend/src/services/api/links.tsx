const baseUrl = process.env.REACT_APP_ENV === "production" ? process.env.PUBLIC_URL : "http://localhost:8080";

export interface LinkData {
    url: string;
    name: string;
    description: string;
    views?: number
}

interface CreateLinkResponse {
    error?: string;
}

export interface GetLinksResponse {
    links: LinkData[];
    error?: string;
}

export const createLink = async (link: LinkData): Promise<CreateLinkResponse> => {
    try {
        const response = await fetch(`${baseUrl}/${link.name}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(link)
        });
        if (!response.ok) {
            let message = 'Failed to complete'
            if (response.body) {
                try {
                    const body = await response.json()
                    message = body.error
                } catch (exception: any) {
                    message = response.statusText
                }
            }
            return { error: message }
        }
        return {}
    } catch (exception: any) {
        return { error: 'Failed to communicate with server, please try your request again' }
    }
}

export const getPopular = async (): Promise<GetLinksResponse> => {
    try {
        const response = await fetch(`${baseUrl}/api/popular`);
        if (!response.ok) {
            let message = 'Failed to complete'
            if (response.body) {
                const body = await response.json()
                message = body.error
            }
            return { links: [], error: message }
        }
        const links: LinkData[] = await response.json()
        return { links }
    } catch (exception: any) {
        return { links: [], error: 'Failed to communicate with server, please try your request again' }
    }
}

export const getRecent = async (): Promise<GetLinksResponse> => {
    try {
        const response = await fetch(`${baseUrl}/api/recent`);
        if (!response.ok) {
            let message = 'Failed to complete'
            if (response.body) {
                const body = await response.json()
                message = body.error
            }
            return { links: [], error: message }
        }
        const links: LinkData[] = await response.json()
        return { links }
    } catch (exception: any) {
        return { links: [], error: 'Failed to communicate with server, please try your request again' }
    }
}

export const getOwned = async (): Promise<GetLinksResponse> => {
    try {
        const response = await fetch(`${baseUrl}/api/owned`);
        if (!response.ok) {
            let message = 'Failed to complete'
            if (response.body) {
                const body = await response.json()
                message = body.error
            }
            return { links: [], error: message }
        }
        const links: LinkData[] = await response.json()
        return { links }
    } catch (exception: any) {
        return { links: [], error: 'Failed to communicate with server, please try your request again' }
    }
}