import './Divider.scss'

export function Divider(props) {
    return <div className={props.secondary ? 'Divider-secondary' : 'Divider-primary'}></div>
}