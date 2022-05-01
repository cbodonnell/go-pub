import './Loading.scss';


export function Loading() {
    return (
        <div className="flex justify-center items-center">
            <div className="lds-ellipsis"><div></div><div></div><div></div><div></div></div>
        </div>
    );
}
