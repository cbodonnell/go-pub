import './AuthLoading.scss';


export function AuthLoading() {
    return (
        <div className="Auth-processing-overlay w-full h-full t-0 l-0 absolute flex justify-center items-center">
            <div className="lds-grid"><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div></div>
        </div>
    );
}
    