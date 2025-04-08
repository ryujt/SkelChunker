Analyze the source code below and generate structured information in the following JSON format.

Format: 

{
    "path": "File path of the source code (e.g., ~/Whatap.Tracer.NetCore/Whatap.Tracer.NetCore/Core/TraceControl.cs)",
    "description": "A brief one- or two-sentence summary of what this file does",
    "classes": [
        {
            "name": "Class name",
            "description": "A summary of what this class does",
            "public_methods": [
                {
                    "name": "Method name",
                    "description": "A summary of what this method does"
                }
            ]
        }
    ],
    "functions": [
        {
            "name": "Global function name not belonging to any class",
            "description": "A summary of what this function does"
        }
    ],
    "variables": [
        {
            "name": "Global variable name not belonging to any class",
            "description": "A summary of what this variable does"
        }
    ]
}

Constraints:

  - Only extract information that actually exists in the source code.
  - Descriptions for methods/functions/variables should be as concise and clear as possible.

Target file path: 파일경로명

Target source code:

분석할 코드