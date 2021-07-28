export function process_example(files) {
	return files
		.map(file => {
			const [name, type] = file.name.split('.');
			return { name, type, source: file.source };
		})
		.sort((a, b) => {
			if (a.name === 'main' && a.type === 'go') return -1;
			if (b.name === 'main' && b.type === 'go') return 1;

			if (a.type === b.type) return a.name < b.name ? -1 : 1;

			if (a.type === 'go') return -1;
			if (b.type === 'go') return 1;
		});
}
