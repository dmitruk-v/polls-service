'use strict';

const App = () => {
	const url = new URL(window.location.href);
	const fields = [];
	for (const [key, value] of url.searchParams) {
		fields.push({ key: key, value: value })
	}
	return (
		<React.Fragment>
			<h1>Survey {SURVEY_ID}</h1>
			<Form fields={fields} />
		</React.Fragment>
	)
}

const Form = ({ fields }) => {
	return (
		<form action="#" className="pt-3">
			{fields.map(item => <FormField key={item.key} label={item.key} value={item.value} />)}
			<button type="submit" className="mt-2 btn btn-primary">Submit</button>
		</form>
	)
}

const FormField = ({ label, value }) => {
	let [input, setInput] = React.useState(value);

	const handleChange = (event) => {
		console.log(event.target.value);
		setInput(event.target.value);
	}

	return (
		<div className="mb-3">
			<label className="form-label">{label}</label>
			<input type="text" value={input} onChange={handleChange} className="form-control" />
		</div>
	)
}

const container = document.getElementById("app");
const root = ReactDOM.createRoot(container);
root.render(<App />);