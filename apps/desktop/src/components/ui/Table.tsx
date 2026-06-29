// SPDX-License-Identifier: MIT
interface Column {
  key: string;
  header: string;
  render?: (value: any, row: any) => React.ReactNode;
  className?: string;
}

interface TableProps {
  columns: Column[];
  data: any[];
  onRowClick?: (row: any) => void;
}

export function Table({ columns, data, onRowClick }: TableProps) {
  if (data.length === 0) {
    return <p className="text-gray-400 text-sm text-center py-8">Nenhum registro.</p>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-200">
            {columns.map((col) => (
              <th key={col.key} className={`text-left text-xs font-medium text-gray-500 uppercase tracking-wider py-3 px-3 ${col.className || ""}`}>
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row, idx) => (
            <tr
              key={row.id || idx}
              className={`border-b border-gray-100 ${onRowClick ? "cursor-pointer hover:bg-gray-50" : ""}`}
              onClick={() => onRowClick?.(row)}
            >
              {columns.map((col) => (
                <td key={col.key} className={`py-3 px-3 text-gray-700 ${col.className || ""}`}>
                  {col.render ? col.render(row[col.key], row) : row[col.key] ?? "-"}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
