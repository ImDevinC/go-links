const baseUrl = process.env.REACT_APP_ENV === "production" ? "https://api.example.com" : "http://localhost:8080";

export interface LinkData {
    url: string;
    name: string;
    description: string;
    views?: number
}

interface CreateLinkResponse {
    error?: string;
}

interface GetLinksResponse {
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
                const body = await response.json()
                message = body.error
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