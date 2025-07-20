import {
  registerDecorator,
  ValidationOptions,
  ValidationArguments,
} from 'class-validator';

export function MaxPathDepth(
  max: number,
  validationOptions?: ValidationOptions,
) {
  return function (object: object, propertyName: string) {
    registerDecorator({
      name: 'maxPathDepth',
      target: object.constructor,
      propertyName,
      options: validationOptions,
      validator: {
        validate(value: string, args: ValidationArguments) {
          return value.split('/').length <= max;
        },
        defaultMessage(args: ValidationArguments) {
          return `path must not exceed ${max} levels`;
        },
      },
    });
  };
}
